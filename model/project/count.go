package project

func (project *Project) Count() (n int, err error) {
	c := project.GetC()
	defer c.Database.Session.Close()

	n, err = c.Find(nil).Count()
	if err != nil {
		err = nil
	}

	return
}
