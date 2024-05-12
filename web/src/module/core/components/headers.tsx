import logo from "../../../assets/logo.svg";
import login from "../../../assets/login.svg";

export const Header = () => {
  return (
    <div className="py-4 flex flex-row justify-between z-10">
      <img src={logo} alt="logo" />
      <button className="flex flex-row bg-[#181E29] text-[16px] text-white px-6 py-2 rounded-3xl border-[#353C4A] border-[1px]">
        Login
        <img src={login} alt="login" className="ml-2 w-4" />
      </button>
    </div>
  );
};
