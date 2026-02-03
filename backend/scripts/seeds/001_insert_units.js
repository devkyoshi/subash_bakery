// Seed: Insert default units and unit conversions
// Date: 2025-12-16
// Description: Populates the units and unit_charts collections with common units of measure
// Note: Units are global/system-wide, not organization-specific

db = db.getSiblingDB('erp_db');

const currentDate = new Date();

print('Starting seed: Insert default units...');

// ============================================
// 1. Insert Units
// ============================================

const units = [
  // Quantity Units
  {
    _id: ObjectId(),
    code: 'PCS',
    name: 'Pieces',
    symbol: 'pcs',
    description: 'Individual pieces or items',
    unit_type: 'quantity',
    is_base_unit: true,
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  {
    _id: ObjectId(),
    code: 'BOX',
    name: 'Box',
    symbol: 'box',
    description: 'Box containing multiple pieces',
    unit_type: 'quantity',
    is_base_unit: false,
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  {
    _id: ObjectId(),
    code: 'CTN',
    name: 'Carton',
    symbol: 'ctn',
    description: 'Carton containing multiple boxes',
    unit_type: 'quantity',
    is_base_unit: false,
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  {
    _id: ObjectId(),
    code: 'PALLET',
    name: 'Pallet',
    symbol: 'plt',
    description: 'Pallet containing multiple cartons',
    unit_type: 'quantity',
    is_base_unit: false,
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  
  // Weight Units
  {
    _id: ObjectId(),
    code: 'G',
    name: 'Gram',
    symbol: 'g',
    description: 'Unit of mass - Gram',
    unit_type: 'weight',
    is_base_unit: true,
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  {
    _id: ObjectId(),
    code: 'KG',
    name: 'Kilogram',
    symbol: 'kg',
    description: 'Unit of mass - Kilogram',
    unit_type: 'weight',
    is_base_unit: false,
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  {
    _id: ObjectId(),
    code: 'LB',
    name: 'Pound',
    symbol: 'lb',
    description: 'Unit of mass - Pound',
    unit_type: 'weight',
    is_base_unit: false,
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  {
    _id: ObjectId(),
    code: 'TON',
    name: 'Metric Ton',
    symbol: 't',
    description: 'Unit of mass - Metric Ton',
    unit_type: 'weight',
    is_base_unit: false,
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  
  // Volume Units
  {
    _id: ObjectId(),
    code: 'ML',
    name: 'Milliliter',
    symbol: 'ml',
    description: 'Unit of volume - Milliliter',
    unit_type: 'volume',
    is_base_unit: true,
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  {
    _id: ObjectId(),
    code: 'L',
    name: 'Liter',
    symbol: 'l',
    description: 'Unit of volume - Liter',
    unit_type: 'volume',
    is_base_unit: false,
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  {
    _id: ObjectId(),
    code: 'GAL',
    name: 'Gallon',
    symbol: 'gal',
    description: 'Unit of volume - Gallon (US)',
    unit_type: 'volume',
    is_base_unit: false,
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  
  // Length Units
  {
    _id: ObjectId(),
    code: 'MM',
    name: 'Millimeter',
    symbol: 'mm',
    description: 'Unit of length - Millimeter',
    unit_type: 'length',
    is_base_unit: true,
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  {
    _id: ObjectId(),
    code: 'CM',
    name: 'Centimeter',
    symbol: 'cm',
    description: 'Unit of length - Centimeter',
    unit_type: 'length',
    is_base_unit: false,
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  {
    _id: ObjectId(),
    code: 'M',
    name: 'Meter',
    symbol: 'm',
    description: 'Unit of length - Meter',
    unit_type: 'length',
    is_base_unit: false,
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  {
    _id: ObjectId(),
    code: 'IN',
    name: 'Inch',
    symbol: 'in',
    description: 'Unit of length - Inch',
    unit_type: 'length',
    is_base_unit: false,
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  {
    _id: ObjectId(),
    code: 'FT',
    name: 'Foot',
    symbol: 'ft',
    description: 'Unit of length - Foot',
    unit_type: 'length',
    is_base_unit: false,
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  }
];

// Store unit IDs for use in conversion charts
const unitMap = {};
units.forEach(unit => {
  unitMap[unit.code] = unit._id;
});

const insertResult = db.units.insertMany(units);
print(`Inserted ${Object.keys(insertResult.insertedIds).length} units`);
print(`Units inserted with IDs: ${Object.values(insertResult.insertedIds).map(id => id.toString()).join(', ')}`);
print('');

// ============================================
// 2. Insert Unit Conversion Charts
// ============================================

const unitCharts = [
  // Quantity conversions
  {
    _id: ObjectId(),
    from_unit_id: unitMap['BOX'],
    to_unit_id: unitMap['PCS'],
    conversion_rate: 12.0, // 1 BOX = 12 PCS
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  {
    _id: ObjectId(),
    from_unit_id: unitMap['CTN'],
    to_unit_id: unitMap['BOX'],
    conversion_rate: 10.0, // 1 CARTON = 10 BOXES
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  {
    _id: ObjectId(),
    from_unit_id: unitMap['PALLET'],
    to_unit_id: unitMap['CTN'],
    conversion_rate: 20.0, // 1 PALLET = 20 CARTONS
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  
  // Weight conversions
  {
    _id: ObjectId(),
    from_unit_id: unitMap['KG'],
    to_unit_id: unitMap['G'],
    conversion_rate: 1000.0, // 1 KG = 1000 G
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  {
    _id: ObjectId(),
    from_unit_id: unitMap['LB'],
    to_unit_id: unitMap['G'],
    conversion_rate: 453.592, // 1 LB = 453.592 G
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  {
    _id: ObjectId(),
    from_unit_id: unitMap['TON'],
    to_unit_id: unitMap['KG'],
    conversion_rate: 1000.0, // 1 TON = 1000 KG
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  
  // Volume conversions
  {
    _id: ObjectId(),
    from_unit_id: unitMap['L'],
    to_unit_id: unitMap['ML'],
    conversion_rate: 1000.0, // 1 L = 1000 ML
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  {
    _id: ObjectId(),
    from_unit_id: unitMap['GAL'],
    to_unit_id: unitMap['ML'],
    conversion_rate: 3785.41, // 1 GAL (US) = 3785.41 ML
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  
  // Length conversions
  {
    _id: ObjectId(),
    from_unit_id: unitMap['CM'],
    to_unit_id: unitMap['MM'],
    conversion_rate: 10.0, // 1 CM = 10 MM
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  {
    _id: ObjectId(),
    from_unit_id: unitMap['M'],
    to_unit_id: unitMap['CM'],
    conversion_rate: 100.0, // 1 M = 100 CM
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  {
    _id: ObjectId(),
    from_unit_id: unitMap['IN'],
    to_unit_id: unitMap['MM'],
    conversion_rate: 25.4, // 1 IN = 25.4 MM
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  },
  {
    _id: ObjectId(),
    from_unit_id: unitMap['FT'],
    to_unit_id: unitMap['IN'],
    conversion_rate: 12.0, // 1 FT = 12 IN
    is_active: true,
    created_at: currentDate,
    updated_at: currentDate,
    version: 1,
    metadata: {}
  }
];

const chartResult = db.unit_charts.insertMany(unitCharts);
print(`Inserted ${Object.keys(chartResult.insertedIds).length} unit conversion charts`);
print('');

// ============================================
// 3. Create Indexes
// ============================================

// Units indexes
db.units.createIndex({ code: 1 }, { unique: true });
db.units.createIndex({ unit_type: 1 });
db.units.createIndex({ is_active: 1 });

// Unit charts indexes
db.unit_charts.createIndex({ from_unit_id: 1, to_unit_id: 1 }, { unique: true });
db.unit_charts.createIndex({ is_active: 1 });

print('Seed completed successfully!');
print('---');
print('Summary:');
print(`- Units created: ${units.length}`);
print(`- Conversion charts created: ${unitCharts.length}`);
print('- Units are global/system-wide (no organization_id)');
