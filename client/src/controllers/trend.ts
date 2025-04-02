import { TrendService } from "../services/trend";
import { TrendSchema, TrendGraphSchema, MonthTrafficTrendSchema } from "../schemas/trend";
import { getPastMonth } from "../utils/date";

export class TrendController {
    /**
     * Get traffic trend of an OLT device.
     * 
     * @param {string} olt OLT device.
     * @returns {TrendGraphSchema[]} Trend of the OLT device.
     */
    static async getTrend(olt: string): Promise<TrendGraphSchema[]> {
        const monthConsult = getPastMonth();
        const response = await TrendService.getTrend(olt, monthConsult);
        if (response.status === 200 && response.data) {
            let trend: TrendSchema = response.data;
            let trends: TrendGraphSchema[] = [];
            trend.months.map((traffic: MonthTrafficTrendSchema) => {
                let currentTraffic: TrendGraphSchema = {
                    month: traffic.month,
                    in: traffic.in,
                    out: traffic.out
                }
                trends.push(currentTraffic);
            });
            let newTrafficTrend: TrendGraphSchema = {
                month: new Date().getMonth() + 1,
                in: trend.out_trend[0].predictions,
                out: trend.out_trend[0].predictions
            }
            trends.push(newTrafficTrend);
            return trends;
        } else {
            console.error(response.err!.message);
            return [];
        }
}
}