export function unit(value: number): string {
    if (typeof value === 'number') {
        if (value >= 1e12) return (value / 1e12).toFixed(2) + " Tb";
        else if (value >= 1e9) return (value / 1e9).toFixed(2) + " Gb";
        else if (value >= 1e6) return (value / 1e6).toFixed(2) + " Mb";
        else if (value >= 1e3) return (value / 1e3).toFixed(2) + " Kb";
        else return value + " Bits";
    } else return value;
}