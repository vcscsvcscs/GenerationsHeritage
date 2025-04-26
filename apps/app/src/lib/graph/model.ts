export interface NodeMenu {
    onClick: () => void;
    deleteNode: () => void;
    createRelationshipAndNode: () => void;
    addRelationship: () => void;
    addRecipe: (() => void) | undefined;
    addAdmin: (() => void) | undefined;
    top: number | undefined;
    left: number | undefined;
    right: number | undefined;
    bottom: number | undefined;
}