import { AxiosError } from "axios";
import axios from "../../core/services/axios.ts";

export const redirectToUrl = async (shortId: string) => {
  try {
    const axiosRes = await axios.get(`/url-shorten/${shortId}`);
    console.log(axiosRes);
    if (axiosRes.status === 302) {
      // Redirect the user to the specified URL
      window.location.href = axiosRes.headers.location;
    } else {
      // Handle other types of responses
    }
    return axiosRes.data;
  } catch (error: AxiosError) {
    throw error.response.data;
  }
};

export const createShotrUrl = async (url: string) => {
  try {
    const axiosRes = await axios.post("/url-shorten", { url });
    return axiosRes.data;
  } catch (error: AxiosError) {
    throw error.response.data;
  }
};
