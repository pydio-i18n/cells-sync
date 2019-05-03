/*
 * Copyright (c) 2019. Abstrium SAS <team (at) pydio.com>
 * This file is part of Pydio Cells.
 *
 * Pydio Cells is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Pydio Cells is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Pydio Cells.  If not, see <http://www.gnu.org/licenses/>.
 *
 * The latest code can be found at <https://pydio.com>.
 */

package control

import (
	"sync"
	"time"

	"github.com/pydio/sync/config"

	"github.com/pydio/cells/common/sync/merger"
)

type SyncStatus int

const (
	SyncStatusIdle SyncStatus = iota
	SyncStatusPaused
	SyncStatusDisabled
	SyncStatusProcessing
	SyncStatusError
)

type EndpointInfo struct {
	FoldersCount   uint64
	FilesCount     uint64
	TotalSpace     uint64
	AvailableSpace uint64
}

type SyncState struct {
	UUID      string
	Config    *config.Task
	Connected bool

	// Sync Process
	Status            SyncStatus
	LastConnection    time.Time
	LastSyncTime      time.Time
	LastProcessStatus merger.ProcessStatus

	// Endpoints Current Info
	LeftInfo  *EndpointInfo
	RightInfo *EndpointInfo
}

type StateStore interface {
	LastState() SyncState
	UpdateConnection(c bool) SyncState
	UpdateSyncStatus(s SyncStatus) SyncState
	UpdateProcessStatus(processStatus merger.ProcessStatus, status ...SyncStatus) SyncState
}

type MemoryStateStore struct {
	sync.Mutex
	config *config.Task
	state  SyncState
}

func NewMemoryStateStore(config *config.Task) *MemoryStateStore {
	s := &MemoryStateStore{
		config: config,
		state: SyncState{
			UUID:      config.Uuid,
			Config:    config,
			Status:    SyncStatusIdle,
			Connected: false,
		},
	}
	return s
}

func (b *MemoryStateStore) LastState() SyncState {
	b.Lock()
	defer b.Unlock()
	return b.state
}

func (b *MemoryStateStore) UpdateSyncStatus(s SyncStatus) SyncState {
	b.Lock()
	defer b.Unlock()
	b.state.Status = s
	return b.state
}

func (b *MemoryStateStore) UpdateProcessStatus(processStatus merger.ProcessStatus, status ...SyncStatus) SyncState {
	b.Lock()
	defer b.Unlock()
	b.state.LastSyncTime = time.Now()
	b.state.LastProcessStatus = processStatus
	if len(status) > 0 {
		b.state.Status = status[0]
	}
	return b.state
}

func (b *MemoryStateStore) UpdateConnection(c bool) SyncState {
	b.Lock()
	defer b.Unlock()
	b.state.Connected = c
	if c {
		b.state.LastConnection = time.Now()
	}
	return b.state
}