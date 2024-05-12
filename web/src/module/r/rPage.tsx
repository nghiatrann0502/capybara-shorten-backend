import { useParams } from "react-router-dom";
import Lottie from "react-lottie";
import loading from "../../assets/loading.json";

const RPage = () => {
  const { id } = useParams();
  if (id)
    window.location.href = `http://localhost:8000/api/v1/url-shorten/${id}`;

  return (
    <div className="w-full h-screen flex justify-center items-center">
      <Lottie
        options={{
          loop: true,
          autoplay: true,
          animationData: loading,
        }}
      />
    </div>
  );
};

export default RPage;
