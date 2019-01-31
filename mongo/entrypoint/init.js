db.createUser(
    {
        user: "msteger",
        pwd: "msteger",
        roles:[
            {
                role: "readWrite",
                db:   "rbtracker"
            }
        ]
    }
);