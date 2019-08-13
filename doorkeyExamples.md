# Doorkey Examples  
  
*// = comments*  
**Bold = Result**  
  
*// Assign variables, create a function, call the function with parameters*  
let five = 5;  
  
let ten = 10;  
  
let add = fn(x, y) {  
  x + y;  
};  
  
let result = add(five, ten);  
  
result  
**15**  
  
*// Create and insert elements into an array*  
let arr = [1,2,3];  
arr  
**[1, 2, 3]**  
  
*// Call builtin len() to count elements in array*  
len(arr)  
**3**  
  
*// Call the first element in an array*  
first(arr)  
**1**  
  
*// Call the last element in an array*  
last(arr)  
**3**  
  
*// Call all elements except the first, in an array*  
tail(arr)  
**[2,3]**  
  
*// Create a function*  
let square = fn(x){x*x};  
  
*// Call the second array element as a parameter within a function*  
square(arr[1])  
**4**  
  
*// Call the builtin push() function to add an element to the array.*  
let arr = push(arr, "four");  
arr  
**[1, 2, 3, four]**  
  
*// Use operators on array elements*  
arr[1] * arr[2]  
**6**  
  
*// Concatenate string element inside array*  
arr[3] + "th"  
**fourth**  
  
*// Concatenate string inside array and push into the array*  
let arr = push(arr, arr[5] + "th");  
arr  
**[1,2,3,four,fourth]**  
  
*// Create a hash table*  
let books = [{"title": "The Sea-Wolf", "authorFirstName": "Jack", "authorLastName": "London", "age": 115}, {"title": "Fear and Loathing in Las Vegas", "authorFirstName": "Hunter", "authorLastName": "Thompson", "age": 48}];  
  
books  
**[{title: The Sea-Wolf, authorFirstName: Jack, authorLastName: London, age: 115}, {title: Fear and Loathing in Las Vegas, authorFirstName: Hunter, authorLastName: Thompson, age: 48}]**  
  
books[0]["title"]  
**The Sea-Wolf**  
  
*// Create a mapping function*  
let a = [1,3,5,7];  
  
let square = fn(x) {x*x};  
  
let map = fn(newArr, f) {let iter = fn(newArr, accumulated) { if (len(newArr) == 0) {accumulated} else{iter(tail(newArr),push(accumulated, f(first(newArr))));}}; iter(newArr, []);};  
  
map(a,square);  
**[1, 9, 25, 49]**  
