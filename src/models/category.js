/**
 * Category Model
 * Represents a product category
 */

export class Category {
  constructor(data) {
    this.id = data.id;
    this.name = data.name;
    this.description = data.description;
    this.parentCategoryId = data.parentCategoryId;
    this.createdAt = data.createdAt;
    this.updatedAt = data.updatedAt;
  }

  toJSON() {
    return {
      id: this.id,
      name: this.name,
      description: this.description,
      parentCategoryId: this.parentCategoryId,
      createdAt: this.createdAt,
      updatedAt: this.updatedAt
    };
  }

  isRootCategory() {
    return !this.parentCategoryId;
  }
}

