import { PayloadAction, createSlice } from "@reduxjs/toolkit";
import type { RootState } from "./store";

// Define a type for the slice state
interface UserState {
  token: string;
}

// Define the initial state using that type
const initialState: UserState = {
  token: "",
};

export const userSlice = createSlice({
  name: "user",
  initialState,
  reducers: {
    set: (state, action: PayloadAction<string>) => {
      state.token = action.payload;
    },
  },
});

export const { set } = userSlice.actions;

// Other code such as selectors can use the imported `RootState` type
export const selectToken = (state: RootState) => state.user.token;

export default userSlice.reducer;

export type { UserState };
