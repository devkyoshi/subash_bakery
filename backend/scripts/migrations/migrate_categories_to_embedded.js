// MongoDB Migration Script: Convert Category Hierarchy to Embedded Structure
// Run this script using: mongosh <database_name> migration_categories.js

// This script migrates categories from parent-child structure to embedded subcategories
// It also updates products to include subcategory_id references

print("=== Starting Category Structure Migration ===\n");

// Step 1: Find all root categories (level 0 or parent_id is null)
print("Step 1: Finding root categories...");
const rootCategories = db.categories.find({
  $or: [
    { level: 0 },
    { parent_id: null },
    { parent_id: { $exists: false } }
  ],
  deleted_at: null
}).toArray();

print(`Found ${rootCategories.length} root categories\n`);

// Step 2: For each root category, find its children and embed them
let migratedCategories = 0;
let migratedSubcategories = 0;
let subcategoriesToDelete = [];

rootCategories.forEach(rootCategory => {
  print(`Processing category: ${rootCategory.name} (${rootCategory._id})`);
  
  // Find all direct children (level 1 or parent_id matches this category)
  const children = db.categories.find({
    $or: [
      { parent_id: rootCategory._id },
      { level: 1, path: new RegExp(`^${rootCategory.path}/`) }
    ],
    deleted_at: null
  }).toArray();
  
  print(`  Found ${children.length} subcategories`);
  
  // Convert children to subcategory format
  const subcategories = children.map(child => {
    // Update products that reference this child category
    const productUpdateResult = db.products.updateMany(
      { category_id: rootCategory._id, deleted_at: null },
      { $set: { subcategory_id: child._id } }
    );
    
    print(`    - ${child.name}: Updated ${productUpdateResult.modifiedCount} products with subcategory_id`);
    migratedSubcategories++;
    
    // Mark this child category for deletion
    subcategoriesToDelete.push(child._id);
    
    return {
      _id: child._id,
      name: child.name,
      code: child.code || "",
      description: child.description || "",
      is_active: child.is_active !== undefined ? child.is_active : true,
      product_count: child.product_count || 0,
      metadata: child.metadata || {},
      created_at: child.created_at,
      updated_at: child.updated_at,
      deleted_at: null,
      created_by: child.created_by,
      updated_by: child.updated_by,
      deleted_by: child.deleted_by,
      version: child.version || 0
    };
  });
  
  // Update the root category with embedded subcategories
  // Remove parent_id, level, and path fields
  db.categories.updateOne(
    { _id: rootCategory._id },
    {
      $set: {
        subcategories: subcategories,
        updated_at: new Date()
      },
      $unset: {
        parent_id: "",
        level: "",
        path: ""
      }
    }
  );
  
  migratedCategories++;
  print(`  ✓ Migrated category with ${subcategories.length} subcategories\n`);
});

// Step 3: Delete the old child category documents
print("\nStep 3: Cleaning up old subcategory documents...");
if (subcategoriesToDelete.length > 0) {
  const deleteResult = db.categories.deleteMany({
    _id: { $in: subcategoriesToDelete }
  });
  print(`Deleted ${deleteResult.deletedCount} old subcategory documents\n`);
} else {
  print("No subcategory documents to delete\n");
}

// Step 4: Remove obsolete fields from all remaining categories
print("Step 4: Removing obsolete fields from categories...");
db.categories.updateMany(
  {},
  {
    $unset: {
      parent_id: "",
      level: "",
      path: ""
    }
  }
);

// Step 5: Ensure all categories have subcategories array (even if empty)
print("Step 5: Ensuring all categories have subcategories array...");
db.categories.updateMany(
  { subcategories: { $exists: false } },
  { $set: { subcategories: [] } }
);

// Step 6: Verify migration
print("\n=== Migration Verification ===");

const totalCategories = db.categories.countDocuments({ deleted_at: null });
print(`Total categories after migration: ${totalCategories}`);

const categoriesWithSubcategories = db.categories.countDocuments({
  deleted_at: null,
  "subcategories.0": { $exists: true }
});
print(`Categories with subcategories: ${categoriesWithSubcategories}`);

const totalSubcategories = db.categories.aggregate([
  { $match: { deleted_at: null } },
  { $project: { subcategoryCount: { $size: "$subcategories" } } },
  { $group: { _id: null, total: { $sum: "$subcategoryCount" } } }
]).toArray();

if (totalSubcategories.length > 0) {
  print(`Total embedded subcategories: ${totalSubcategories[0].total}`);
}

const productsWithSubcategory = db.products.countDocuments({
  deleted_at: null,
  subcategory_id: { $exists: true, $ne: null }
});
print(`Products with subcategory_id: ${productsWithSubcategory}`);

// Step 7: Create indexes for better performance
print("\n=== Creating Indexes ===");

// Index for category lookups
db.categories.createIndex({ organization_id: 1, deleted_at: 1 });
db.categories.createIndex({ organization_id: 1, name: 1 }, { unique: true, partialFilterExpression: { deleted_at: null } });
db.categories.createIndex({ "subcategories._id": 1 });

// Index for product queries by category and subcategory
db.products.createIndex({ category_id: 1, subcategory_id: 1, deleted_at: 1 });
db.products.createIndex({ organization_id: 1, category_id: 1, deleted_at: 1 });

print("Indexes created successfully");

print("\n=== Migration Summary ===");
print(`✓ Migrated ${migratedCategories} root categories`);
print(`✓ Embedded ${migratedSubcategories} subcategories`);
print(`✓ Updated products with subcategory references`);
print(`✓ Removed obsolete fields (parent_id, level, path)`);
print(`✓ Created necessary indexes`);
print("\n=== Migration Complete ===");

// Optional: Rollback script
print("\n--- Rollback Information ---");
print("To rollback this migration, you would need to:");
print("1. Extract subcategories from embedded arrays and create separate documents");
print("2. Restore parent_id, level, and path fields");
print("3. Remove subcategory_id from products");
print("4. This migration creates a backup of the original structure is recommended before running.");
