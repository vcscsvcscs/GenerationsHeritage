import { Node, Relationship, Date } from 'neo4j-driver';
import { Integer } from 'neo4j-driver';

export interface PersonProperties {
    allow_admin_access: boolean;
    google_id: string;
    first_name: string;
    middle_name?: string;
    last_name: string;
    titles?: string[]; // e.g. Jr., Sr., III
    suffixes?: string[]; // e.g. Ph.D., M.D.
    extra_names?: string[];
    aliases?: string[];
    mothers_first_name?: string;
    mothers_last_name?: string;
    born?: Date<number>;
    place_of_birth?: string;
    died?: Date<number>;
    place_of_death?: string;
    life_events?: { [key: string]: {from: Date<number>, to:Date<number>, desription: string} }[];
    occupations?: string[];
    occupation_to_display?: string;
    others_said?: { [key: string]: string };
    limit: number;
    photos?: { [key: string]: string };
    videos?: { [key: string]: string };
    audios?: { [key: string]: string };
    profile_picture?: string;
    verified: boolean;
    email?: string;
    phone?: string;
    residence?: {
        city?: string;
        country?: string;
        zip_code?: string;
        address_line_1?: string;
        address_line_2?: string;
    };
    religion?: string;
    baptized?: string;
    ideology?: string;
    blood_type?: string;
    allergies?: string[];
    medications?: string[];
    medical_conditions?: string[];
    height?: number;
    weight?: number;
    hair_colour?: string;
    skin_colour?: string;
    eye_colour?: string;
    sports?: string[];
    hobbies?: string[];
    interests?: string[];
    languages?: { [key: string]: string };
    notes?: string;
}

export interface FamilyRelationship {
    verified: boolean;
    notes?: string;
    from?: Date<number>;
    to?: Date<number>;
}

export type Person = Node<Integer, PersonProperties>;
export type Parent = Relationship<Integer, FamilyRelationship>;
export type Sibling = Relationship<Integer, FamilyRelationship>;
export type Spouse = Relationship<Integer, FamilyRelationship>;
export type Child = Relationship<Integer, FamilyRelationship>;

export interface RecipeProperties {
    id: string;
    name: string;
    origin: string;
    category: string;
    first_recorded: Date<number>;
    description: string;
    ingredients: string[];
    instructions: string[];
    photo: string;
    notes?: string;
}

export type Recipe = Node<Integer, RecipeProperties>;
export type Likes = Relationship<Integer, {
    favourite: boolean;
    like_it: boolean;
    could_make_it: boolean;
}>;

export interface FamilyTree {
    ancestors: String;
    prel1: FamilyRelationship;
    children: String;
    prel2: FamilyRelationship;
    spouses: String;
    srel: FamilyRelationship;
    user: String;

}

export interface FamilyMember {
    person: Person;
}