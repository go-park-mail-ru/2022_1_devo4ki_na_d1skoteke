package app

import (
	"context"
	"cotion/internal/api/domain/entity"
	"cotion/internal/api/domain/repository"
	pb "cotion/internal/note/infra/grpc"
	"cotion/internal/pkg/generator"
	"errors"
	log "github.com/sirupsen/logrus"
)

const packageName = "app notes"

var ErrNoteAccess = errors.New("The user does not have access to this note. Or the note does not exist.")

type NotesApp struct {
	notesRepository      repository.NotesRepository
	usersNotesRepository repository.UsersNotesRepository
	pb.UnimplementedNoteServiceServer
}

func NewNotesApp(notesRepo repository.NotesRepository, usersNotesRepository repository.UsersNotesRepository) *NotesApp {
	return &NotesApp{
		notesRepository:      notesRepo,
		usersNotesRepository: usersNotesRepository,
	}
}

func (n *NotesApp) NotesList(noteReq *pb.NoteReq, stream pb.NoteService_NotesListServer) error {
	notes, err := n.usersNotesRepository.AllNotesByUserID(noteReq.UserID)
	if err != nil {
		return err
	}
	for _, note := range notes {
		grpcNote := &pb.Note{
			Name: note.Name,
			Body: note.Body,
		}
		if err := stream.Send(grpcNote); err != nil {
			return err
		}
	}
	return nil
}

func (n *NotesApp) Save(ctx context.Context, noteReq *pb.NoteReq) (*pb.Nothing, error) {
	logger := log.WithFields(log.Fields{
		"package":  packageName,
		"function": "SaveNote",
	})

	newToken := generator.RandToken()
	newNote := entity.Note{
		Name: noteReq.Note.Name,
		Body: noteReq.Note.Body,
	}

	if err := n.notesRepository.Save(newToken, newNote); err != nil {
		logger.Error(err)
		return &pb.Nothing{Status: false}, err
	}

	if err := n.usersNotesRepository.AddLink(noteReq.UserID, newToken); err != nil {
		logger.Error(err)
		return &pb.Nothing{Status: false}, err
	}

	return &pb.Nothing{Status: true}, nil
}

func (n *NotesApp) Get(ctx context.Context, NoteReq *pb.NoteReq) (*pb.Note, error) {
	logger := log.WithFields(log.Fields{
		"package":  packageName,
		"function": "GetNote",
	})

	if !n.usersNotesRepository.CheckLink(NoteReq.UserID, NoteReq.NoteID) {
		logger.Warning(ErrNoteAccess)
		return nil, ErrNoteAccess
	}

	note, err := n.notesRepository.Find(NoteReq.NoteID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	grpcNote := &pb.Note{
		Name: note.Name,
		Body: note.Body,
	}
	return grpcNote, nil
}

func (n *NotesApp) Update(ctx context.Context, NoteReq *pb.NoteReq) (*pb.Nothing, error) {
	if !n.usersNotesRepository.CheckLink(NoteReq.UserID, NoteReq.NoteID) {
		log.WithFields(log.Fields{
			"package":  packageName,
			"function": "UpdateNote",
		}).Warning(ErrNoteAccess)
		return &pb.Nothing{Status: true}, ErrNoteAccess
	}

	updateNote := entity.Note{
		Name: NoteReq.Note.Name,
		Body: NoteReq.Note.Body,
	}
	if err := n.notesRepository.Update(NoteReq.NoteID, updateNote); err != nil {
		return &pb.Nothing{Status: true}, err
	}
	return &pb.Nothing{Status: true}, nil
}

func (n *NotesApp) Delete(ctx context.Context, NoteReq *pb.NoteReq) (*pb.Nothing, error) {
	logger := log.WithFields(log.Fields{
		"package":  packageName,
		"function": "DeleteNote",
	})

	if !n.usersNotesRepository.CheckLink(NoteReq.UserID, NoteReq.NoteID) {
		logger.Warning(ErrNoteAccess)
		return &pb.Nothing{Status: true}, ErrNoteAccess
	}

	if err := n.notesRepository.Delete(NoteReq.NoteID); err != nil {
		logger.Error(err)
		return &pb.Nothing{Status: true}, err
	}

	return &pb.Nothing{Status: true}, n.usersNotesRepository.DeleteLink(NoteReq.UserID, NoteReq.NoteID)
}
