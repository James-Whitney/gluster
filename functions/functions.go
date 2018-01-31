package functions

type DistFunc int

type Blah struct{
	A int
	B int
}

func (t *DistFunc) Multiply(args *Blah, reply *int) error {
	*reply = args.A * args.B
	return nil
}