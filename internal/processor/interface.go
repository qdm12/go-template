package processor

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Interface

var _ Interface = (*Processor)(nil)

type Interface interface {
	UserCreator
	UserGetter
}
