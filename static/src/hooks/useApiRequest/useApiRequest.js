import {useReducer, useCallback, useContext} from "react";
import axios from "axios";
import reducer, { initialState } from "./reducer";
import { fetching, success, error } from "./actionCreators";
import {AuthContext} from "../../contexts/AuthContext";

const useApiRequest = (endpoint, { verb = "get", params = {}, data = {}} = {}) => {
  const [state, dispatch] = useReducer(reducer, initialState);

  const {isAuthenticated, accessToken, refreshToken} = useContext(AuthContext);

  const makeRequest = useCallback(async () => {
    dispatch(fetching());
    try {
      const headers = {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      };
      if(isAuthenticated){
        console.log("Authenticated!!!", accessToken)
        headers["Authorization"] = `Bearer ${accessToken}`
      }
      console.log("headers", {headers : headers})
      console.log("data", {data : data})
      const url = 'http://127.0.0.1:8000' + endpoint
      const response = await axios({method: verb, url: url, data: data, params: params, headers : headers});
      dispatch(success(response));
    } catch (e) {
      e.data = e.response.data
      console.log(e.response.data)
      dispatch(error(e));
    }
  }, [endpoint, verb, params]);

  return [state, makeRequest];
};

export default useApiRequest;
