import { Node, Relationship, Date } from 'neo4j-driver';
import { Integer } from 'neo4j-driver';

interface PersonProperties {
    id: string;
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
    born?: Date;
    place_of_birth?: string;
    died?: Date;
    place_of_death?: string;
    life_events?: { [key: string]: {from: Date, to:Date, desription: string} }[];
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

export type Person = Node<Integer, PersonProperties>;
export type FamilyRelationship = Relationship<Integer, {
    type: string;
    verified: boolean;
    notes?: string;
    from?: Date;
    to?: Date;
}>;


interface RecipeProperties {
    id: string;
    name: string;
    origin: string;
    category: string;
    first_recorded: Date;
    description: string;
    ingredients: string[];
    instructions: string[];
    photo: string;
    notes?: string;
}

export type Recipe = Node<Integer, RecipeProperties>;
export type RecipeRelationship = Relationship<Integer, {
    favourite: boolean;
    like_it: boolean;
    could_make_it: boolean;
}>;



