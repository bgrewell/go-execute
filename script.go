package execute

type ScriptType string

const (
	ScriptTypePowerShell ScriptType = "powershell"
	ScriptTypeBash       ScriptType = "bash"
	ScriptTypePython     ScriptType = "python"
)
