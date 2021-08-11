package server

import (
	"context"
	"errors"
	"github.com/shlande/dmhy-rss/pkg/controller/manager"
	"github.com/shlande/dmhy-rss/pkg/controller/manual"
	"github.com/shlande/dmhy-rss/pkg/controller/mission"
	"github.com/shlande/dmhy-rss/pkg/controller/store"
	"github.com/shlande/dmhy-rss/pkg/data"
	"github.com/shlande/dmhy-rss/pkg/data/provider"
	"github.com/shlande/dmhy-rss/pkg/data/tools"
	"github.com/sirupsen/logrus"
)

type Server struct {
	store store.Interface

	manager *manager.Manager

	manual *manual.Manual

	pvd provider.Provider

	clpd *tools.CollectionProvider
}

func (s *Server) Run(ctx context.Context) {
	go s.manager.Run(ctx)
}

func (s *Server) SetAnimationKeywords(ctx context.Context, anmId, keywords string) error {
	anm, err := s.store.Animation().Get(anmId)
	if err != nil {
		return err
	}
	err = anm.SetKeywords(keywords)
	if err != nil {
		return err
	}
	return s.store.Animation().Save(anm)
}

func (s *Server) ListCollectionByAnmId(ctx context.Context, anmUniId string) ([]*data.Collection, error) {
	anm, err := s.pvd.Get(ctx, anmUniId)
	if err != nil {
		return nil, err
	}
	return s.manual.Search(ctx, anm)
}

func (s *Server) ListCollectionByKeywords(ctx context.Context, keywords string) ([]*data.Collection, error) {
	return s.clpd.Keywords(ctx, keywords)
}

func (s *Server) SearchAnimation(ctx context.Context, keywords string) ([]*data.Animation, error) {
	return s.pvd.Search(ctx, keywords)
}

func (s *Server) GetAnimationBySession(ctx context.Context, year int, session int) ([]*data.Animation, error) {
	var ss provider.Session
	switch session {
	case 0:
		ss = provider.Spring
	case 1:
		ss = provider.Summer
	case 2:
		ss = provider.Fall
	case 4:
		ss = provider.Winter
	default:
		return nil, errors.New("无效的session")
	}
	return s.pvd.Session(ctx, year, ss)
}

func (s *Server) CreateMission(_ context.Context, collectionId string) error {
	return s.manual.CreateMission(collectionId)
}

func (s *Server) CancelMission(ctx context.Context, collectionId string) error {
	return errors.New("这个功能还没有去完成哦")
}

func (s *Server) ListMission(_ context.Context, active bool) (mss []*mission.Mission, err error) {
	if active {
		for _, v := range s.manager.ListActiveMissionId() {
			ms, err := s.store.Mission().Get(v)
			if err != nil {
				logrus.Error(err)
				continue
			}
			mss = append(mss, ms)
		}
		return mss, err
	} else {
		return s.store.Mission().GetAll(false)
	}
}

func (s *Server) GetLogs(_ context.Context, missionId string) ([]*mission.Log, error) {
	return s.store.Log().GetAll(missionId)
}

func (s *Server) GetCollection(ctx context.Context, collectionId string) (*data.Collection, error) {
	return s.manual.Get(ctx, collectionId)
}
