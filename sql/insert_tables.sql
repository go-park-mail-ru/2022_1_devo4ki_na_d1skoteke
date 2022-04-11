INSERT INTO CotionUser(UserID, Username, Email, Password)
VALUES (1, 'Test', 'test@mail.ru', 'test');

INSERT INTO Note(NoteID, Name, Body)
VALUES ('1', '1st psql note', 'Body of 1st psql note.'),
       ('3', '3st psql note', 'Body of 3st psql note.');

INSERT INTO UsersNotes(UserID, NoteID)
VALUES (1, 1);