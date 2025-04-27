import type { components } from '$lib/api/api.gen.js';

export function validatePersonRegistration(data: components['schemas']['PersonRegistration']): string | null {
    if (!data.first_name || data.first_name.trim() === "") {
        return "First name is required.";
    }

    if (!data.last_name || data.last_name.trim() === "") {
        return "Last name is required.";
    }

    if (data.email !== undefined && data.email !== null && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(data.email)) {
        return "Invalid email format.";
    }

    if (!data.born || isNaN(Date.parse(data.born))) {
        return "Valid birth date is required.";
    }

    if (!data.biological_sex || !['male', 'female', 'intersex', 'unknown', 'other'].includes(data.biological_sex.toString())){
        return 'Invalid value for biological sex. Must be male female, intersex, unknown, or other.';
    }

    if (!data.mothers_first_name || data.mothers_first_name.trim() === "") {
        return "Mother's first name is required.";
    }

    if (!data.mothers_last_name || data.mothers_last_name.trim() === "") {
        return "Mother's last name is required.";
    }

    return null; // No errors
}

export function validateFamilyRelationship(relationship: components['schemas']['FamilyRelationship'] & {type:string}): string | null {
    const validRelationships = [
        "child",
        "parent",
        "spouse",
        "sibling"
    ];

    if (!validRelationships.includes(relationship.type)) {
        return `Invalid family relationship. Must be one of ${validRelationships.join(', ')}.`;
    }

    if (relationship.from !== undefined && relationship.from !== null && isNaN(Date.parse(relationship.from))) {
        return "Valid date is required for 'from' field.";
    }

    if (relationship.to !== undefined && relationship.to !== null && isNaN(Date.parse(relationship.to))) {
        return "Valid date is required for 'to' field.";
    }

    return null; // No errors
}