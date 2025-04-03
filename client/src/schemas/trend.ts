export interface PredictionTrendSchema {
    mse: number;
    predictions: number;
}

export interface MonthTrafficTrendSchema {
    month: number;
    in: number;
    out: number;
}

export interface TrendSchema {
    months: MonthTrafficTrendSchema[];
    out_trend: PredictionTrendSchema[];
    in_trend: PredictionTrendSchema[];
}

export interface TrendGraphSchema {
    month: number;
    in: number;
    out: number;
}