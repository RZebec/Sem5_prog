package ticket

const firstTestTicket = `{
"Info": {
"Id": 1,
"Title": "TestTitle",
"Editor": {
"Mail": "",
"UserId": 0,
"FirstName": "",
"LastName": ""
},
"HasEditor": false,
"Creator": {
"Mail": "test@test.de",
"FirstName": "Max",
"LastName": "Mustermann"
},
"CreationTime": "2018-10-27T18:23:04.0276343+02:00",
"LastModificationTime": "2018-10-27T18:23:04.0276343+02:00"
},
"Messages": [
{
"Id": 0,
"CreatorMail": "test@test.de",
"CreationTime": "2018-09-27T18:23:04.0831377+02:00",
"Content": "TestContent1",
"OnlyInternal": false
}
]
}`

const secondTestTicket = `{
    "Info": {
        "Id": 2,
        "Title": "TestTitle2",
        "Editor": {
            "Mail": "",
            "UserId": 0,
            "FirstName": "",
            "LastName": ""
        },
        "HasEditor": false,
        "Creator": {
            "Mail": "test@test.de",
            "FirstName": "Max",
            "LastName": "Mustermann"
        },
        "CreationTime": "2018-10-27T18:23:04.0896374+02:00",
        "LastModificationTime": "2018-10-27T18:23:04.0896374+02:00"
    },
    "Messages": [
        {
            "Id": 0,
            "CreatorMail": "test@test.de",
            "CreationTime": "2018-10-27T18:23:04.0951371+02:00",
            "Content": "TestContent2",
            "OnlyInternal": false
        }
    ]
}`

const thirdTestTicket = `{
    "Info": {
        "Id": 3,
        "Title": "TestTitle",
        "Editor": {
            "Mail": "test@test.de",
            "UserId": 1,
            "FirstName": "Max",
            "LastName": "Mustermann"
        },
        "HasEditor": true,
        "Creator": {
            "Mail": "test@test.de",
            "FirstName": "Max",
            "LastName": "Mustermann"
        },
        "CreationTime": "2018-10-27T18:23:04.1016353+02:00",
        "LastModificationTime": "2018-10-27T18:23:04.1016353+02:00"
    },
    "Messages": [
        {
            "Id": 0,
            "CreatorMail": "test@test.de",
            "CreationTime": "2018-10-27T18:23:04.107137+02:00",
            "Content": "TestContent3",
            "OnlyInternal": false
        }
    ]
}`

const fourthTestTicket = `{
    "Info": {
        "Id": 4,
        "Title": "TestTitle4",
        "Editor": {
            "Mail": "peter@test.de",
            "UserId": 2,
            "FirstName": "Peter",
            "LastName": "Test"
        },
        "HasEditor": true,
        "Creator": {
            "Mail": "peter@test.de",
            "FirstName": "Peter",
            "LastName": "Test"
        },
        "CreationTime": "2018-10-27T18:23:04.1141357+02:00",
        "LastModificationTime": "2018-10-27T18:23:04.1141357+02:00"
    },
    "Messages": [
        {
            "Id": 0,
            "CreatorMail": "peter@test.de",
            "CreationTime": "2018-10-27T18:23:04.1201396+02:00",
            "Content": "TestContent4",
            "OnlyInternal": true
        }
    ]
}`
