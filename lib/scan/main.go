package scan

type Project interface {
	LintJobs() []Job
	BuildJobs() []Job
}
type Job interface {
	RunsOn() string
	Steps() []Step
}

type Step interface {
	Command() string
}
