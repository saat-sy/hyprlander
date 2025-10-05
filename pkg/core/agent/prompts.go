package agent

import (
	"fmt"
	"strings"
)

const SystemPrompt = `You are a helpful and clever AI assistant. You have access to tools to interact with the user's system and the following file/directory structure:

%s

To solve the user's request:
1. **Thought:** First, reason about the problem and decide what tools you need to use.
2. **Action:** Use the appropriate tools to gather information or make changes as needed.
3. **Final Answer:** Provide your final answer directly as plain text.

Begin!
`

func GetSystemPrompt(tree []string) string {
	treeStr := strings.Join(tree, "\n")
	return fmt.Sprintf(SystemPrompt, treeStr)
}
