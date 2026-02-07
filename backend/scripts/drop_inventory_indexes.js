// Drop all indexes for inventory-related collections
db = db.getSiblingDB('erp_db');

const collections = [
    'units',
    'unit_charts',
    'stock_levels',
    'stock_movements',
    'batches',
    'stock_adjustments',
    'inventory_counts',
    'serial_numbers'
];

collections.forEach(colName => {
    print(`Processing collection: ${colName}`);
    try {
        // dropIndexes() removes all indexes except the default _id index
        const result = db[colName].dropIndexes();
        print(`  Result: ${JSON.stringify(result)}`);
    } catch (e) {
        // Often errors if the collection doesn't exist or has no extra indexes
        print(`  Note: ${e.message}`);
    }
});

print('Finished dropping indexes.');
