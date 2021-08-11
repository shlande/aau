package remote

import (
	"github.com/shlande/dmhy-rss/pkg/data"
)

type Server struct {
	UnimplementedSubscriberServer
	addPeer    map[chan *AddedResponse]struct{}
	createPeer map[chan *CreatedResponse]struct{}
}

func (s Server) created(collection *data.Collection) {
	cl := newCollection(collection)
	for c, _ := range s.createPeer {
		c <- &CreatedResponse{Collection: cl}
	}
}

func (s Server) added(source *data.Source, collection *data.Collection) {
	cl := newCollection(collection)
	sc := newResource(source)
	for c, _ := range s.addPeer {
		c <- &AddedResponse{Collection: cl, Resource: sc}
	}
}

func (s Server) Created(collection *data.Collection) {
	go s.created(collection)
}

func (s Server) Added(detail *data.Source, collection *data.Collection) {
	go s.added(detail, collection)
}

func (s Server) AddedStream(request *ConsumeRequest, server Subscriber_AddedStreamServer) error {
	ch := make(chan *AddedResponse)
	s.addPeer[ch] = struct{}{}
	defer delete(s.addPeer, ch)

	for resp := range ch {
		if err := server.Send(resp); err != nil {
			return err
		}
	}
	return nil
}

func (s Server) CreatedStream(request *ConsumeRequest, server Subscriber_CreatedStreamServer) error {
	ch := make(chan *CreatedResponse)
	s.createPeer[ch] = struct{}{}
	defer delete(s.createPeer, ch)

	for resp := range ch {
		if err := server.Send(resp); err != nil {
			return err
		}
	}
	return nil
}

func newCollection(collection *data.Collection) *Collection {
	return &Collection{
		Animation: newAnimation(collection.Animation),
		Latest:    int64(collection.Latest),
	}
}

func newAnimation(animation *data.Animation) *Animation {
	return &Animation{
		Id:           animation.Id,
		Name:         animation.Name,
		Translated:   animation.Translated,
		AirDate:      animation.AirDate.Unix(),
		TotalEpisode: int64(animation.TotalEpisodes),
		Category:     animation.Category,
	}
}

func newResource(source *data.Source) *Resource {
	return &Resource{
		Name:        source.Name,
		CreateTime:  source.CreateTime.Unix(),
		DownloadUrl: source.GetDownloadUrl(),
		Episode:     int64(source.Episode),
		Metadata:    newMetadata(&source.Metadata),
	}
}

func newMetadata(md *data.Metadata) *Metadata {
	return &Metadata{
		Fansub:   md.Fansub,
		Quality:  int64(md.Quality),
		Language: int64(md.Language),
		Subtype:  int64(md.SubType),
		Type:     "",
	}
}
