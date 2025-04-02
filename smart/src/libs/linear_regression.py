import pandas as pd
from sklearn.linear_model import LinearRegression
from sklearn.metrics import mean_squared_error
from sklearn.model_selection import train_test_split
from src.utils.change import bits_to_gbits, gbits_to_bits


class RegressionLineal:
    @staticmethod
    def get_out(df_interface, month: int):
        X_train, X_test, y_train, y_test = train_test_split(
            df_interface[["Month"]],
            df_interface[["Out"]],
            test_size=0.2,
            random_state=123,
        )

        model = LinearRegression()
        model.fit(X_train, y_train)

        y_pred = model.predict(X_test)

        mse = mean_squared_error(y_test, y_pred)

        new_month = pd.DataFrame({"Month": [month]})
        predictions = model.predict(new_month)
        return mse, predictions

    @staticmethod
    def get_in(df_interface, month: int):
        X_train, X_test, y_train, y_test = train_test_split(
            df_interface[["Month"]],
            df_interface[["In"]],
            test_size=0.2,
            random_state=123,
        )

        model = LinearRegression()
        model.fit(X_train, y_train)

        y_pred = model.predict(X_test)

        mse = mean_squared_error(y_test, y_pred)

        new_month = pd.DataFrame({"Month": [month]})
        predictions = model.predict(new_month)
        return mse, predictions

    @staticmethod
    def run_procress(trends: list, month: int):
        data = []
        for trend in trends:
            data.append(
                {
                    "Device": trend.device_id,
                    "Month": trend.date.month,
                    "Out": bits_to_gbits(trend.out),
                    "In": bits_to_gbits(trend.in_),
                    "Bandwidth": bits_to_gbits(trend.bandwidth),
                }
            )

        df = pd.DataFrame(data)
        out = []
        in_ = []

        mseIn, predictionsIn = RegressionLineal.get_in(df, month)
        in_.append(
            {"mse": mseIn, "predictions": gbits_to_bits(float(predictionsIn[0][0]))}
        )
        mseOut, predictionsOut = RegressionLineal.get_out(df, month)
        out.append(
            {"mse": mseOut, "predictions": gbits_to_bits(float(predictionsOut[0][0]))}
        )

        return out, in_
