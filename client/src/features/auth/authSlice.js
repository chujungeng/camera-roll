import { createSlice } from '@reduxjs/toolkit'

export const authSlice = createSlice({
    name: "auth",
    initialState: {
        loggedIn: false,
    },
    reducers: {
        toggleLogin: (state) => {
            state.loggedIn = !state.loggedIn
        },

        logOut: (state) => {
            state.loggedIn = false
        },

        logIn: (state) => {
            state.loggedIn = true
        },
    }
});

export const { toggleLogin, logOut, logIn } = authSlice.actions;

export default authSlice.reducer;