INSERT INTO CotionUser(UserID, Username, Email, Password)
VALUES ('c04532ca4e12438bcd37d2ae1676d3f5a27241062095eaccdbf0102b78d2a948', 'Test account', 'test@mail.ru', '758ab49634fc498f25a7149f4cfb2b9594ddd962a9f3a4546125004a5ebebe61'),
       ('b0e251ff2bb51b963aed043ba5a92c867e9f4a5ab0bb6906fc5bfe26b932e1d7', 'Nikita account', 'nikita@mail.ru', '6796c96c7d3190e4ae976e038858cc294e773eac43d867dd5eb151d59e02349b');

INSERT INTO Note(NoteID, Name, Body)
VALUES ('1', '1st psql note', 'Body of 1st psql note.'),
       ('3', '3st psql note', 'Body of 3st psql note.');

INSERT INTO UsersNotes(UserID, NoteID)
VALUES ('c04532ca4e12438bcd37d2ae1676d3f5a27241062095eaccdbf0102b78d2a948', '1'),
       ('c04532ca4e12438bcd37d2ae1676d3f5a27241062095eaccdbf0102b78d2a948', '3');