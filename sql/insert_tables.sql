INSERT INTO CotionUser(UserID, Username, Email, Password)
VALUES (1, 'Test', 'test@mail.ru', 'test');

INSERT INTO Note(NoteID, Name, Body)
VALUES (1, 'name', 'body');

INSERT INTO UsersNotes(UserID, NoteID)
VALUES (1, 1);