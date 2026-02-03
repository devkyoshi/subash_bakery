// Fix timestamps stored in Extended JSON format ({ "$date": "..." }) to native BSON Dates

print("Starting timestamp fix migration...");

var count = 0;

db.users.find({}).forEach(function(doc) {
    var updates = {};
    var dirty = false;

    // Helper to fix a specific date field
    function fixDateField(fieldName) {
        var val = doc[fieldName];
        // Check if field exists and has the $date structure
        if (val && val.$date) {
            var dateVal;
            
            // Handle ISO string format: { "$date": "2023-01-01T..." }
            if (typeof val.$date === 'string') {
                dateVal = new Date(val.$date);
            } 
            // Handle NumberLong format: { "$date": { "$numberLong": "..." } }
            else if (typeof val.$date === 'object' && val.$date.$numberLong) {
                dateVal = new Date(parseInt(val.$date.$numberLong));
            }
            
            if (dateVal) {
                updates[fieldName] = dateVal;
                dirty = true;
            }
        }
    }

    // List of date fields to check and fix
    var dateFields = [
        'created_at', 
        'updated_at', 
        'last_login', 
        'last_activity', 
        'email_verified_at', 
        'phone_verified_at',
        'password_changed_at',
        'last_failed_login_at',
        'locked_until',
        'deleted_at'
    ];

    dateFields.forEach(function(field) {
        fixDateField(field);
    });

    if (dirty) {
        db.users.updateOne({ _id: doc._id }, { $set: updates });
        print("Fixed timestamps for user: " + doc._id);
        count++;
    }
});

print("Migration completed. Updated " + count + " documents.");
