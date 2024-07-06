import { configureStore } from "@reduxjs/toolkit";
import UserReducer from "./userReducer";
import NavReducer from "./navReducer";
export const store = configureStore({
  reducer: {
    user: UserReducer,
    nav: NavReducer
  },
});
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
