package nougat

// Methods

// Head sets the Nougat method to HEAD and sets the given pathURL.
func (r *Nougat) Head(pathURL string) *Nougat {
	r.method = "HEAD"
	return r.Path(pathURL)
}

// Get sets the Nougat method to GET and sets the given pathURL.
func (r *Nougat) Get(pathURL string) *Nougat {
	r.method = "GET"
	return r.Path(pathURL)
}

// Post sets the Nougat method to POST and sets the given pathURL.
func (r *Nougat) Post(pathURL string) *Nougat {
	r.method = "POST"
	return r.Path(pathURL)
}

// Put sets the Nougat method to PUT and sets the given pathURL.
func (r *Nougat) Put(pathURL string) *Nougat {
	r.method = "PUT"
	return r.Path(pathURL)
}

// Patch sets the Nougat method to PATCH and sets the given pathURL.
func (r *Nougat) Patch(pathURL string) *Nougat {
	r.method = "PATCH"
	return r.Path(pathURL)
}

// Delete sets the Nougat method to DELETE and sets the given pathURL.
func (r *Nougat) Delete(pathURL string) *Nougat {
	r.method = "DELETE"
	return r.Path(pathURL)
}

// Options sets the Nougat method to OPTIONS and sets the given pathURL.
func (r *Nougat) Options(pathURL string) *Nougat {
	r.method = "OPTIONS"
	return r.Path(pathURL)
}

// Trace sets the Nougat method to TRACE and sets the given pathURL.
func (r *Nougat) Trace(pathURL string) *Nougat {
	r.method = "TRACE"
	return r.Path(pathURL)
}

// Connect sets the Nougat method to CONNECT and sets the given pathURL.
func (r *Nougat) Connect(pathURL string) *Nougat {
	r.method = "CONNECT"
	return r.Path(pathURL)
}
