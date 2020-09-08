package fileobserver

import (
	"context"
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/reactivex/rxgo/observable"
	"github.com/reactivex/rxgo/observer"
	"path/filepath"
)


func StartMonitor(observingPath string, ctx context.Context) (err error, ob observable.Observable) {
	if observingPath == "" {
		err = errors.New("observing path is not provided")
		return
	}

	ob = observable.Create(func(emitter *observer.Observer, disposed bool) {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			emitter.OnError(err)
			return
		}
		defer watcher.Close()

		go func() {
			if ctx.Done() != nil {
				<-ctx.Done()
				emitter.OnDone()
				watcher.Close()
			}
		}()

		err = watcher.Add(observingPath)
		if err != nil {
			emitter.OnError(err)
		}

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op & fsnotify.Create == fsnotify.Create {
					s, err := filepath.Abs(event.Name)
					if err != nil {
						emitter.OnError(err)
						return
					}
					emitter.OnNext(s)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				emitter.OnError(err)
			}
		}
	})

	return
}
