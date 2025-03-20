import { LoadStatus } from "../constant/loadStatus";

/**
 * Declaration of all types of cargo status.
 */
export type LoadingStateValue = LoadStatus.EMPTY | LoadStatus.LOADING | LoadStatus.LOADED;