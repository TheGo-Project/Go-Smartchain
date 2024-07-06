import { PayloadAction, createSlice } from "@reduxjs/toolkit";
import type { RootState } from "./store";

// Define a type for the slice state
interface NavState {
  active: string;
}

// Define the initial state using that type
const initialState: NavState = {
  active: "home",
};

export const navSlice = createSlice({
  name: "nav",
  initialState,
  reducers: {
    set: (state, action: PayloadAction<string>) => {
      state.active = action.payload;
    },
  },
});

export const { set } = navSlice.actions;

// Other code such as selectors can use the imported `RootState` type
export const selectActive = (state: RootState) => state.nav.active;

export default navSlice.reducer;

export type { NavState };
