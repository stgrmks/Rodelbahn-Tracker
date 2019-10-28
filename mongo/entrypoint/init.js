db.auth('admin', 'admin');
db.createUser(
    {
        user: "user",
        pwd: "user",
        roles:[
            {
                role: "readWrite",
                db:   "rb-tracker"
            }
        ]
    }
);