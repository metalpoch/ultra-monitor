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

        modelo = LinearRegression()
        modelo.fit(X_train, y_train)

        y_pred = modelo.predict(X_test)

        mse = mean_squared_error(y_test, y_pred)

        nuevo_mes = pd.DataFrame({"Month": [month]})
        predicciones = modelo.predict(nuevo_mes)
        return mse, predicciones

    @staticmethod
    def get_in(df_interface, month: int):
        X_train, X_test, y_train, y_test = train_test_split(
            df_interface[["Month"]],
            df_interface[["In"]],
            test_size=0.2,
            random_state=123,
        )

        modelo = LinearRegression()
        modelo.fit(X_train, y_train)

        y_pred = modelo.predict(X_test)

        mse = mean_squared_error(y_test, y_pred)

        nuevo_mes = pd.DataFrame({"Month": [month]})
        predicciones = modelo.predict(nuevo_mes)
        return mse, predicciones

    @staticmethod
    def run_procress(trends: list, month: int):
        data = []
        for trend in trends:
            data.append(
                {
                    "Device": trend.Device_Id,
                    "Month": trend.date.month,
                    "Out": bits_to_gbits(trend.Out),
                    "In": bits_to_gbits(trend.In),
                    "Bandwidth": bits_to_gbits(trend.Bandwidth),
                }
            )

        df = pd.DataFrame(data)
        out = []
        in_ = []

        mseIn, prediccionIn = RegressionLineal.get_in(df, month)
        in_.append(
            {"mse": mseIn, "prediccion": gbits_to_bits(float(prediccionIn[0][0]))}
        )
        mseOut, prediccionOut = RegressionLineal.get_out(df, month)
        out.append(
            {"mse": mseOut, "prediccion": gbits_to_bits(float(prediccionOut[0][0]))}
        )

        return out, in_
