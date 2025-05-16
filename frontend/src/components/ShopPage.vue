<template>
  <div class="shop-page">
    <h4>惠农商城</h4>
    <!-- 商品列表 -->
    <div class="products-grid">
      <div v-for="(product, index) in products" 
           :key="index" 
           class="product-card">
        <img :src="product.image" :alt="product.name">
        <h4>{{product.name}}</h4>
        <p class="price">¥{{product.price}}</p>
        <button class="buy-btn">立即购买</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ShopPage',
  data() {
    return {
      products: [
        {
          name: '优质大米',
          price: '29.9',
          image: require('@/assets/fcb741694431e0ca5984d74892dad2fa.png')
        },
        {
          name: '新鲜蔬菜',
          price: '15.9',
          image: require('@/assets/8c860028b7c2d04ec3ef93d04e2725ed.png')
        }
      ]
    }
  }
}
</script>

<style scoped>
.shop-page {
  padding: 15px;
  max-width: 800px;
  margin: 0 auto;
}

.shop-page h4 {
  text-align: center;
  margin-bottom: 20px;
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 15px;
  max-width: 600px;
  margin: 0 auto;
}

.product-card {
  background: #fff;
  border-radius: 8px;
  padding: 10px;
  text-align: center;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.product-card img {
  width: 100%;
  height: 120px;
  object-fit: cover;
  border-radius: 4px;
}

.price {
  color: #f56c6c;
  font-weight: bold;
}

.buy-btn {
  background: #2C5530;
  color: white;
  border: none;
  padding: 5px 15px;
  border-radius: 15px;
  margin-top: 5px;
  cursor: pointer;
}

h4 {
margin: 0px;
}

</style>
