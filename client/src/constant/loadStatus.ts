/**
 * @enum Status of system loading.
 * 
 * - EMPTY: Status when no search has been done.
 * - LOADING: Status when a search is being done.
 * - LOADED: Status when a search has been done.
 */
export enum LoadStatus {
    EMPTY = "EMPTY",
    LOADING = "LOADING",
    LOADED = "LOADED"
}