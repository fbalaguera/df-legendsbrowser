package model

import (
	"strings"

	"github.com/robertjanetzko/LegendsBrowser2/backend/util"
	"golang.org/x/exp/slices"
)

func (w *DfWorld) process() {

	// set site in structure
	for _, site := range w.Sites {
		for _, structure := range site.Structures {
			structure.SiteId = site.Id_
		}
	}

	w.processEvents()
	w.processCollections()

	// check events texts
	for _, e := range w.HistoricalEvents {
		e.Details.Html(&Context{World: w})
	}
}

func (w *DfWorld) processEvents() {
	for _, e := range w.HistoricalEvents {
		switch d := e.Details.(type) {
		case *HistoricalEventHfDoesInteraction:
			if strings.HasPrefix(d.Interaction, "DEITY_CURSE_WEREBEAST_") {
				w.HistoricalFigures[d.TargetHfid].Werebeast = true
			}
			if strings.HasPrefix(d.Interaction, "DEITY_CURSE_VAMPIRE_") {
				w.HistoricalFigures[d.TargetHfid].Vampire = true
			}
		case *HistoricalEventCreatedSite:
			w.addEntitySite(d.CivId, d.SiteId)
			w.addEntitySite(d.SiteCivId, d.SiteId)
		case *HistoricalEventDestroyedSite:
			w.addEntitySite(d.DefenderCivId, d.SiteId)
			w.addEntitySite(d.SiteCivId, d.SiteId)
			w.Sites[d.SiteId].Ruin = true
		case *HistoricalEventSiteTakenOver:
			w.addEntitySite(d.AttackerCivId, d.SiteId)
			w.addEntitySite(d.SiteCivId, d.SiteId)
			w.addEntitySite(d.DefenderCivId, d.SiteId)
			w.addEntitySite(d.NewSiteCivId, d.SiteId)
		case *HistoricalEventHfDestroyedSite:
			w.addEntitySite(d.SiteCivId, d.SiteId)
			w.addEntitySite(d.DefenderCivId, d.SiteId)
			w.Sites[d.SiteId].Ruin = true
		case *HistoricalEventReclaimSite:
			w.addEntitySite(d.SiteCivId, d.SiteId)
			w.addEntitySite(d.SiteCivId, d.SiteId)
			w.Sites[d.SiteId].Ruin = false
		}
	}
}

func (w *DfWorld) processCollections() {
	for _, col := range w.HistoricalEventCollections {
		for _, eventId := range col.Event {
			if e, ok := w.HistoricalEvents[eventId]; ok {
				e.Collection = col.Id_
			}
		}

		switch cd := col.Details.(type) {
		case *HistoricalEventCollectionAbduction:
			targets := make(map[int]bool)
			for _, eventId := range col.Event {
				if e, ok := w.HistoricalEvents[eventId]; ok {
					switch d := e.Details.(type) {
					case *HistoricalEventHfAbducted:
						targets[d.TargetHfid] = true
					}
				}
			}
			delete(targets, -1)
			cd.TargetHfids = util.Keys(targets)
		case *HistoricalEventCollectionBeastAttack:
			attackers := make(map[int]bool)
			for _, eventId := range col.Event {
				if e, ok := w.HistoricalEvents[eventId]; ok {
					switch d := e.Details.(type) {
					case *HistoricalEventHfSimpleBattleEvent:
						attackers[d.Group1Hfid] = true
					case *HistoricalEventHfAttackedSite:
						attackers[d.AttackerHfid] = true
					case *HistoricalEventHfDestroyedSite:
						attackers[d.AttackerHfid] = true
					case *HistoricalEventAddHfEntityLink:
						attackers[d.Hfid] = true
					case *HistoricalEventCreatureDevoured:
						attackers[d.Eater] = true
					case *HistoricalEventItemStolen:
						attackers[d.Histfig] = true
					}
				}
			}
			delete(attackers, -1)
			cd.AttackerHfIds = util.Keys(attackers)
		case *HistoricalEventCollectionJourney:
		HistoricalEventCollectionJourneyLoop:
			for _, eventId := range col.Event {
				if e, ok := w.HistoricalEvents[eventId]; ok {
					switch d := e.Details.(type) {
					case *HistoricalEventHfTravel:
						cd.TravellerHfIds = d.GroupHfid
						break HistoricalEventCollectionJourneyLoop
					}
				}
			}
		case *HistoricalEventCollectionOccasion:
			for _, eventcolId := range col.Eventcol {
				if e, ok := w.HistoricalEventCollections[eventcolId]; ok {
					switch d := e.Details.(type) {
					case *HistoricalEventCollectionCeremony:
						d.OccasionEventcol = col.Id_
					case *HistoricalEventCollectionCompetition:
						d.OccasionEventcol = col.Id_
					case *HistoricalEventCollectionPerformance:
						d.OccasionEventcol = col.Id_
					case *HistoricalEventCollectionProcession:
						d.OccasionEventcol = col.Id_
					}
				}
			}
		}
	}
}

func (w *DfWorld) addEntitySite(entityId, siteId int) {
	if e, ok := w.Entities[entityId]; ok {
		if !slices.Contains(e.Sites, siteId) {
			e.Sites = append(e.Sites, siteId)
		}
	}
}
