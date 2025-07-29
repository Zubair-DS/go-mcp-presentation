04 Dynamic tool discovery
----------------------
Problem: different models support different amount of tools. So therefore 128 tools might be ok with one model but might not be ok for another mode, how to fix?

dynamic tool discovery, it looks at the prompt and decides how many tools to use

DEMO
User settings: "github.copilot.chat.virtualTools.enabled": true
BEFORE: error: too many tools
AFTER: it works
