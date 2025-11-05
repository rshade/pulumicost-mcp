import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";

// Configuration
const config = new pulumi.Config();
const instanceType = config.get("instanceType") || "t3.micro";
const dbInstanceClass = config.get("dbInstanceClass") || "db.t3.small";

// VPC
const vpc = new aws.ec2.Vpc("example-vpc", {
    cidrBlock: "10.0.0.0/16",
    enableDnsHostnames: true,
    enableDnsSupport: true,
    tags: {
        Name: "pulumicost-example-vpc",
        Environment: "dev",
        ManagedBy: "pulumi",
    },
});

// Subnets
const publicSubnet = new aws.ec2.Subnet("public-subnet", {
    vpcId: vpc.id,
    cidrBlock: "10.0.1.0/24",
    availabilityZone: "us-east-1a",
    mapPublicIpOnLaunch: true,
    tags: {
        Name: "pulumicost-example-public-subnet",
        Environment: "dev",
    },
});

const privateSubnet = new aws.ec2.Subnet("private-subnet", {
    vpcId: vpc.id,
    cidrBlock: "10.0.2.0/24",
    availabilityZone: "us-east-1b",
    tags: {
        Name: "pulumicost-example-private-subnet",
        Environment: "dev",
    },
});

// Internet Gateway
const igw = new aws.ec2.InternetGateway("internet-gateway", {
    vpcId: vpc.id,
    tags: {
        Name: "pulumicost-example-igw",
    },
});

// Security Group for EC2
const webSecurityGroup = new aws.ec2.SecurityGroup("web-sg", {
    vpcId: vpc.id,
    description: "Allow HTTP and SSH traffic",
    ingress: [
        { protocol: "tcp", fromPort: 80, toPort: 80, cidrBlocks: ["0.0.0.0/0"] },
        { protocol: "tcp", fromPort: 443, toPort: 443, cidrBlocks: ["0.0.0.0/0"] },
        { protocol: "tcp", fromPort: 22, toPort: 22, cidrBlocks: ["0.0.0.0/0"] },
    ],
    egress: [
        { protocol: "-1", fromPort: 0, toPort: 0, cidrBlocks: ["0.0.0.0/0"] },
    ],
    tags: {
        Name: "pulumicost-example-web-sg",
    },
});

// EC2 Instance
const ami = aws.ec2.getAmi({
    mostRecent: true,
    owners: ["amazon"],
    filters: [
        { name: "name", values: ["amzn2-ami-hvm-*-x86_64-gp2"] },
    ],
});

const webServer = new aws.ec2.Instance("web-server", {
    instanceType: instanceType,
    ami: ami.then(a => a.id),
    subnetId: publicSubnet.id,
    vpcSecurityGroupIds: [webSecurityGroup.id],
    tags: {
        Name: "pulumicost-example-web-server",
        Environment: "dev",
        Purpose: "web",
    },
});

// S3 Bucket
const bucket = new aws.s3.Bucket("app-bucket", {
    acl: "private",
    versioning: {
        enabled: true,
    },
    lifecycleRules: [
        {
            enabled: true,
            transitions: [
                {
                    days: 30,
                    storageClass: "STANDARD_IA",
                },
                {
                    days: 90,
                    storageClass: "GLACIER",
                },
            ],
        },
    ],
    tags: {
        Name: "pulumicost-example-bucket",
        Environment: "dev",
    },
});

// RDS Security Group
const dbSecurityGroup = new aws.ec2.SecurityGroup("db-sg", {
    vpcId: vpc.id,
    description: "Allow PostgreSQL traffic from web servers",
    ingress: [
        {
            protocol: "tcp",
            fromPort: 5432,
            toPort: 5432,
            securityGroups: [webSecurityGroup.id],
        },
    ],
    tags: {
        Name: "pulumicost-example-db-sg",
    },
});

// RDS Subnet Group
const dbSubnetGroup = new aws.rds.SubnetGroup("db-subnet-group", {
    subnetIds: [publicSubnet.id, privateSubnet.id],
    tags: {
        Name: "pulumicost-example-db-subnet-group",
    },
});

// RDS Instance
const database = new aws.rds.Instance("postgres-db", {
    allocatedStorage: 20,
    engine: "postgres",
    engineVersion: "14.7",
    instanceClass: dbInstanceClass,
    dbName: "exampledb",
    username: "dbadmin",
    password: config.requireSecret("dbPassword"),
    skipFinalSnapshot: true,
    vpcSecurityGroupIds: [dbSecurityGroup.id],
    dbSubnetGroupName: dbSubnetGroup.name,
    tags: {
        Name: "pulumicost-example-db",
        Environment: "dev",
    },
});

// Outputs
export const vpcId = vpc.id;
export const publicSubnetId = publicSubnet.id;
export const webServerPublicIp = webServer.publicIp;
export const webServerPublicDns = webServer.publicDns;
export const bucketName = bucket.id;
export const databaseEndpoint = database.endpoint;
