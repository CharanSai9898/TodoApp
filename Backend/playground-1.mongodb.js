/* global use, db */

use('todo_app'); 

// // Delete existing collection
// db.sales.drop();

// // Insert documents without specifying _id
// db.sales.insertMany([
//   { item: 'Apple', price: 10, quantity: 5, category: 'Fruit' },
//   { item: 'Banana', price: 5, quantity: 10, category: 'Fruit' },
//   { item: 'Carrot', price: 8, quantity: 6, category: 'Vegetable' },
//   { item: 'Tomato', price: 6, quantity: 8, category: 'Vegetable' },
//   { item: 'Mango', price: 15, quantity: 3, category: 'Fruit' }
// ]);

// // Display documents
// db.sales.find();

//db.sales.aggregate([{$match:{category:"Fruit"}}])
//db.sales.find({},{_id:0,item:1,price:1,quantity:1})

// db.sales.aggregate([{
//     $project: {
//       _id: 0,
//       item: 1,
//       quantity:1,
//       totalCost:{$multiply:["$price","$quantity"]}
//     }
// }])

// db.sales.aggregate([
//     {
//         $match:{
//             category:'Vegetable'
//         }
//     },
//     {
//         $group:{
//             _id:'$category',
//             totalVegetables:{$sum: 1},
//             totalCost:{
//                 $sum:{
//                     $multiply:['$price','$quantity']
//                 }
//             }
//         }
//     }
// ])

// db.sales.aggregate([
//     {
//         $match:{
//             category:'Fruit'
//         }
//     },
//     {
//         $group:{
//             _id:'$category',
//             totalFruits:{$sum: 1},
//             totalCost:{
//                 $sum:{
//                     $multiply:['$price','$quantity']
//                 }
//             }
//         }
//     }
// ])

// db.sales.aggregate([{
//     $group: {
//       _id: "$category",
//       totalSales: {
//         $sum: {$multiply:["$price","$quantity"]}
//       }
//     }
// }])

// db.sales.aggregate([{
// $match: { category: "Fruit" }
// },{
// $project: {
// item: 1,
// totalCost: {
// $multiply: ["$price", "$quantity"]
// }}},
// {
// $sort: { totalCost: -1 }
// },{
// $limit: 2
// }])