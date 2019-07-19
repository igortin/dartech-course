package abc

type playlist struct {
	node *song
}

type song struct {
	group  string
	name   string
	prev   *song
	nex    *song
}

func (p *playlist) nshow() string {
	if p.node.nex != nil {
		p.node = p.node.nex
		return p.node.name
	}
	return "NO next"
}

func (p *playlist) pshow() string {
	if p.node.prev != nil {
		p.node = p.node.prev
		return p.node.name
	}
	return "No prev"
}

func (p *playlist) append(val *song) {
	if p.node == nil {
		p.node = val
		return
	}
	for s := p.node; s != nil; s = s.nex {
		if s.nex == nil {
			tmp := s
			s.nex = val
			s = s.nex
			s.prev = tmp
		}
	}
}

// Конструтор
func CreateSong(name string, group string, p *song, n *song) song {
	s := song{name, group, nil, nil}
	return s
}

func GetName(s *song) (string, string) {
	return s.name,s.group
}

func Point(p *playlist, s *song){
	p.node = s
}

func App(p *playlist, s *song){
	p.append(s)
}

func GetN(p *playlist) string {
	return p.nshow()
}

func GetP(p *playlist) string {
	return p.pshow()
}

func Instance() playlist{
	return playlist{}
}