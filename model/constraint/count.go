package constraint

func (constraint *Constraint) Count() (n int, err error) {
	c := constraint.GetC()
	defer c.Database.Session.Close()

	n, err = c.Find(nil).Count()
	if err != nil {
		err = nil
	}

	return
}
