package purge


// Job represents a directory cleanup task
type Job struct {
    Dir string
    Cfg *Config
}

// NewJob creates a new purge job
func NewJob(dir string, cfg *Config) *Job {
    return &Job{
        Dir: dir,
        Cfg: cfg,
    }
}

// Plan runs a dry run and returns all changes that would be made
func (j *Job) Plan() ([]Change, error) {
    return PreviewChanges(j.Dir, j.Cfg)
}