interface TitleBarProps {
    title: string;
}

export default function TitleBar({ title }: TitleBarProps) {
    return (
        <div className="w-fit h-fit p-3 flex flex-col justify-center items-center">
            <h1 className="text-2xl font-semibold text-blue-800 text-center">{title}</h1>
            <div className="w-full h-1 bg-blue-800 rounded-full"></div>
        </div>
    )
}