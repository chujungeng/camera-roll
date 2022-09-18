import { createSlice } from '@reduxjs/toolkit'

export const apiSlice = createSlice({
    name: "api",
    initialState: {
        root: process.env.REACT_APP_API_SERVER? process.env.REACT_APP_API_SERVER: "http://localhost:9648/api/admin/",
    },
    reducers: {
        update: (state, action) => {
            state.root = action.payload
        },
    }
});

export const { update } = apiSlice.actions;

export default apiSlice.reducer;