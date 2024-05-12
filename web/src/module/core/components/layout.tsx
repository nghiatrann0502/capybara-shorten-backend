import bg from "../../../assets/bg.svg";
import swirl from "../../../assets/swirl.svg";
import { Header } from "./headers";

interface IMainLayoutProps {
  children: React.ReactNode;
}

export const MainLayout: React.FC<IMainLayoutProps> = ({ children }) => {
  return (
    <div className="bg-[#0B101B] min-h-screen min-w-full relative ">
      <img src={swirl} alt="BG" className="absolute h-full" />
      <img src={bg} alt="BG" className="absolute top-24 md:top-0" />
      <div
        className="h-full px-4 container md:mx-auto"
        style={{ position: "inherit" }}
      >
        <Header />
        {children}
      </div>
      {/* <div className="px-4 container md:mx-auto absolute bottom-0 h-5"> */}
      {/*   <SubFooter /> */}
      {/* </div> */}
    </div>
  );
};
