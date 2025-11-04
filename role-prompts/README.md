# Role-Specific Prompts for PulumiCost MCP Development

This directory contains specialized prompts tailored for different roles working with the PulumiCost MCP server. These prompts help AI assistants provide context-aware guidance based on your current role and objectives.

## Available Roles

- **[Senior Architect](senior-architect.md)**: System design, architecture decisions, scalability
- **[Product Manager](product-manager.md)**: Feature planning, prioritization, user stories
- **[DevOps Engineer](devops-engineer.md)**: Deployment, operations, monitoring
- **[Plugin Developer](plugin-developer.md)**: Plugin creation, spec compliance
- **[Cost Analyst](cost-analyst.md)**: Cost optimization, reporting, budgeting

## How to Use

### With Claude Code

Load a role prompt at the start of your session:

```bash
# In Claude Code, reference the prompt file
@role-prompts/senior-architect.md
```

### In Chat Applications

Copy the relevant prompt and paste it at the beginning of your conversation:

```
I'm working as a [ROLE] on the PulumiCost MCP server project.
[Paste prompt content here]

Now, help me with...
```

### Combining Prompts

You can combine multiple role perspectives:

```
I need both architectural and product perspectives on this feature.
@role-prompts/senior-architect.md
@role-prompts/product-manager.md
```

## Customizing Prompts

Feel free to customize these prompts for your team's specific needs:

1. Copy the base prompt
2. Add team-specific context
3. Include project conventions
4. Reference your documentation

## Contributing New Roles

To add a new role prompt:

1. Create a new `.md` file in this directory
2. Follow the template structure from existing prompts
3. Include role-specific context, responsibilities, and guidelines
4. Update this README with the new role
5. Submit a PR

## Template Structure

```markdown
# [Role Name] - PulumiCost MCP

## Role Context
[Description of the role and responsibilities]

## Key Responsibilities
- [Responsibility 1]
- [Responsibility 2]

## Project Context
[Specific project information relevant to this role]

## Guidelines
[Role-specific guidelines and best practices]

## Common Tasks
[List of common tasks for this role]

## Decision Framework
[How this role makes decisions]
```
