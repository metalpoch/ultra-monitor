import ChangeIcon from "../icons/ChangeIcon"

export default function ChangeStatusDisableButton({id, isDisable, handler}) {
  return <button
    className={`rounded-full p-2 duration-150 hover:bg-slate-950 ${isDisable ? "bg-red-900" : "bg-green-900" }`}
    onClick={() => handler(id, isDisable)}
    >
      <ChangeIcon width={12} height={12} />
  </button>
}
